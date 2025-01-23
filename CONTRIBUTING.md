# Contributing Guidelines

Thank you for your interest in contributing! We welcome and appreciate all contributions, whether they are bug reports, new features, or documentation improvements. This guide will help you get started with contributing to the project.


## Filing Issues

If you encounter a bug, have a feature request, or want to provide general feedback, please file an issue on the GitHub issue tracker for this repository. Before creating a new issue, check the existing issues to avoid duplication. When filing a new issue, include:

- **A descriptive title** and a summary of the issue.
- **Steps to reproduce** the bug (if applicable).
- Relevant **logs**, **screenshots**, or **code snippets** that illustrate the problem.
- Any **suggested solutions** or workarounds (if applicable).

Providing detailed information helps us address your issue promptly and effectively.


## How to Contribute

We’re excited to have you join the project! Follow these steps to contribute:

### 1. Propose Your Change
   - **File an Issue**: Start by submitting an issue that describes your proposed change. This helps us discuss the idea and ensures alignment before any code is written.
   - **Discuss and Get Feedback**: A maintainer will review your issue and may provide suggestions or ask questions. If your proposal is accepted, proceed to the next steps.

### 2. Fork and Clone the Repository
   - **Fork the Repository**:  
     Click the **"Fork"** button at the top-right of the repository page on GitHub to create your own copy of the project under your GitHub username.
     
   - **Clone Your Fork**:  
     Clone your forked repository to your local machine using the following command. Replace `<YourUsername>` with your GitHub username:
     ```bash
     git clone https://github.com/<YourUsername>/csm-sso-plugin.git
     ```
     This command copies the repository to your local machine, allowing you to make changes.
     
     After cloning, navigate to the project directory:
     ```bash
     cd csm-sso-plugin
     ```

### 3. Make Your Changes
   - **Create a New Branch**:  
     It’s important to always create a new branch for each new change or feature. This helps keep your changes organized and avoids conflicts with the `main` branch. To create a new branch, use:
     ```bash
     git checkout -b your-branch-name
     ```
     Replace `your-branch-name` with a descriptive name for the new branch, such as `fix-bug` or `add-new-feature`.

   - **Make Your Changes**:  
     Now that you're on your new branch, make the necessary changes to the code, documentation, or other parts of the project. After making the changes, don’t forget to **stage** and **commit** them:
     ```bash
     git add .
     git commit -m "Your commit message describing the changes"
     ```

     For example, a good commit message might be:
     ```
     git commit -m "Fix issue with user login validation"
     ```

### 4. Run Tests (Optional but Recommended)
   - If the project includes automated tests, run them locally to ensure your changes don’t introduce bugs. Instructions for running tests should be in the project documentation.

### 5. Submit a Pull Request (PR)
   - **Push Your Changes**:
     ```bash
     git push origin your-branch-name
     ```
   - **Open a Pull Request**:  
     Go to the original repository on GitHub and click on **New Pull Request**. Select your branch as the source and describe your changes in detail.
   - **Address Feedback**:  
     A maintainer may review your PR and provide feedback. Be prepared to make changes as requested.


## Tips for a Smooth Contribution Process

- **Keep Commits Atomic**:  
   Ensure each commit represents a single, logical change. This makes it easier to review and understand your contribution.

- **Write Clear Commit Messages**:  
   When writing commit messages, follow the [Conventional Commits](https://www.conventionalcommits.org/) standard (if applicable) for clarity.

- **Stay Synced with the Original Repository**:  
   If the original repository has updates (especially when working on a long-lived branch), you might want to keep your fork up-to-date. You can add the original repository as a remote (called `upstream`) and pull in changes:
   ```bash
   git remote add upstream https://github.com/ThalesGroup/csm-sso-plugin.git
   git fetch upstream
   git merge upstream/main