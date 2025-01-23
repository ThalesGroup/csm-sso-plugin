/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

// Define the fetchRequest function with retry logic
const fetchRequest = async (url, method, data, retries = 10, delay = 5000) => {
  const options = {
    method: method,
    headers: { 'Content-Type': 'application/json' }
  };
  if (data) {
    options.body = JSON.stringify(data);
  }

  for (let attempt = 0; attempt < retries; attempt++) {
    try {
      const response = await fetch(url, options);

      if (!response.ok) {
        const text = await response.text();
        const error = { status: response.status, statusText: response.statusText, text };
        if (response.status === 500) {
          error.statusText = "Internal Server Error";
          error.text = "Unable to reach CipherTrust Manager or Akeyless";
        }
        throw error;
      }

      return await response.json();
    } catch (error) {
      if (attempt < retries - 1) {
        console.log(`Retry attempt ${attempt + 1} for ${url}`);
        await new Promise(resolve => setTimeout(resolve, delay));
      } else {
        if (error instanceof TypeError && error.message === "Failed to fetch") {
          throw {
            status: 503,
            statusText: "Service Unavailable",
            text: "Unable to connect to CSM Executable"
          };
        } else {
          throw {
            status: error.status || 500,
            statusText: error.statusText || "Unknown Error",
            text: error.text || error.message || "An unexpected error has occurred"
          };
        }
      }
    }
  }
};

// Function to get form data
const getFormData = () => {
  return {
    cmUrl: document.getElementById('cmUrl').value.trim(),
    username: document.getElementById('username').value.trim(),
    password: document.getElementById('password').value.trim(),
    apiAccessID: document.getElementById('apiAccessID').value.trim(),
    apiAccessKey: document.getElementById('apiAccessKey').value.trim(),
    akeylessURL: document.getElementById('akeylessURL').value.trim(),
    connectionName: document.getElementById('connectionName').value.trim()
  };
};

// Function to save form data to local storage
const saveFormDataToLocalStorage = () => {
  const formData = getFormData();
  localStorage.setItem('formData', JSON.stringify(formData));
  console.log('Form data saved to local storage:', formData);
};

// Function to load form data from local storage
const loadFormDataFromLocalStorage = () => {
  const storedData = localStorage.getItem('formData');
  console.log('Stored data retrieved from local storage:', storedData);
  if (storedData) {
    try {
      const formData = JSON.parse(storedData);
      document.getElementById('cmUrl').value = formData.cmUrl;
      document.getElementById('username').value = formData.username;
      document.getElementById('password').value = formData.password;
      document.getElementById('apiAccessID').value = formData.apiAccessID;
      document.getElementById('apiAccessKey').value = formData.apiAccessKey;
      document.getElementById('akeylessURL').value = formData.akeylessURL;
      document.getElementById('connectionName').value = formData.connectionName;
    } catch (e) {
      console.error('Failed to parse form data from local storage:', e);
    }
  }
};

// Function to show messages
const showMessage = (message, isError = false) => {
  const messageElement = document.getElementById('message');
  messageElement.textContent = message;
  messageElement.style.backgroundColor = isError ? '#f8d7da' : '#d4edda'; // Error vs success background color
  messageElement.style.color = isError ? '#721c24' : '#155724'; // Error vs success text color
  messageElement.style.display = 'block'; // Show message container
  messageElement.style.border = isError ? '1px solid #cf266e' : '1px solid #369f69';
};

// Function to hide the message container
const hideMessage = () => {
  const messageElement = document.getElementById('message');
  messageElement.style.display = 'none'; // Hide message container
};

// Handle form submission
document.getElementById('configForm').addEventListener('submit', function (event) {
  event.preventDefault();

  // Hide any existing message
  hideMessage();

  const loadingOverlay = document.getElementById('loadingOverlay');

  // Show loading spinner
  loadingOverlay.classList.add('active');

  // Capture form data
  const formData = getFormData();

  // Execute requests in sequence with retry logic
  fetchRequest('http://localhost:52920/initialize', 'POST', {
    cmUrl: formData.cmUrl,
    username: formData.username,
    password: formData.password
  })
    .then(InitializeCMResponse => {
      console.log('Login to CM:', InitializeCMResponse);
      return fetchRequest('http://localhost:52920/csm-service', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(enableCSMResponse => {
      console.log('Check if CSM is enabled:', enableCSMResponse);
      return fetchRequest('http://localhost:52920/create-connection', 'POST', {
        cmUrl: formData.cmUrl,
        connection_name: formData.connectionName,
        akeyless_id: formData.apiAccessID,
        akeyless_key: formData.apiAccessKey,
        akeyless_url: formData.akeylessURL
      });
    })
    .then(createConnectionResponse => {
      console.log('Create Akeyless connection:', createConnectionResponse);
      return fetchRequest('http://localhost:52920/fetch-jwks', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(fetchJWKsResponse => {
      console.log('Fetch JWKs JSON from CM:', fetchJWKsResponse);

      // Adding a 5-second delay before executing the auth-akeyless request
      return new Promise(resolve => setTimeout(resolve, 5000)); // Wait for 5 seconds
    })
    .then(() => {
      return fetchRequest('http://localhost:52920/auth-akeyless', 'POST', {
        cmUrl: formData.cmUrl,
        akeyless_id: formData.apiAccessID,
        akeyless_key: formData.apiAccessKey
      });
    })
    .then(authAkeylessResponse => {
      console.log('Get t-token from Akeyless:', authAkeylessResponse);
      return fetchRequest('http://localhost:52920/create-jwt-auth', 'POST', {
        cmUrl: formData.cmUrl,
        akeyless_id: formData.apiAccessID,
        akeyless_key: formData.apiAccessKey
      });
    })
    .then(createJWTAuthResponse => {
      console.log('Create JWT Auth Method in Akeyless:', createJWTAuthResponse);
      return fetchRequest('http://localhost:52920/create-access-role', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(createAccessRoleResponse => {
      console.log('Create Access Role:', createAccessRoleResponse);
      return fetchRequest('http://localhost:52920/set-role-rule', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(setRulesResponse => {
      console.log('Set rules for Access Role:', setRulesResponse);
      return fetchRequest('http://localhost:52920/associate-role', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(associateAccessRoleResponse => {
      console.log('Associate Access Role:', associateAccessRoleResponse);
      return fetchRequest('http://localhost:52920/update-config', 'POST', {
        cmUrl: formData.cmUrl,
        connection_name: formData.connectionName,
        akeyless_url: formData.akeylessURL
      });
    })
    .then(successResponse => {
      console.log('Update AkeylessConfig:', successResponse);
      showMessage('CSM has been configured successfully!', false);
      return fetchRequest('http://localhost:52920/delete-token', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(deleteTokenResponse => {
      console.log('Delete token response:', deleteTokenResponse);
    })
    .catch(error => {
      console.error('Error:', error);
      showMessage(`${error.text || ''}`, true);
    })
    .finally(() => {
      // Hide loading spinner
      loadingOverlay.classList.remove('active');
    });

  // Save form data to local storage
  saveFormDataToLocalStorage();
});

// Add event listener for the status check button
document.getElementById('checkStatusButton').addEventListener('click', function () {
  // Hide any existing message
  hideMessage();

  const loadingOverlay = document.getElementById('loadingOverlay');

  // Show loading spinner
  loadingOverlay.classList.add('active');

  // Capture form data
  const formData = getFormData();

  // Initialize before fetching status
  fetchRequest('http://localhost:52920/initialize', 'POST', {
    cmUrl: formData.cmUrl,
    username: formData.username,
    password: formData.password
  })
    .then(InitializeCMResponse => {
      console.log('Login to CM:', InitializeCMResponse);
      return fetchRequest('http://localhost:52920/check-status', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(statusCheckResponse => {
      console.log('Check CSM Tile Status:', statusCheckResponse);
      showMessage(statusCheckResponse.message || 'Status check successful!', false);
      return fetchRequest('http://localhost:52920/delete-token', 'POST', {
        cmUrl: formData.cmUrl
      });
    })
    .then(deleteTokenResponse => {
      console.log('Delete token response:', deleteTokenResponse);
    })
    .catch(error => {
      console.error('Error:', error);
      showMessage(`${error.text || ''}`, true);
    })
    .finally(() => {
      // Hide loading spinner
      loadingOverlay.classList.remove('active');
    });
});

// Load form data from local storage on page load
document.addEventListener('DOMContentLoaded', (event) => {
  loadFormDataFromLocalStorage();
});