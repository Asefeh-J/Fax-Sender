# FaxSender - Cross-Platform Fax Sending Application

FaxSender is a cross-platform desktop application built with **Go** for the backend and the **Fyne** package for the user interface. This application allows users to send faxes directly from their desktop. It supports integration with an installer built using **NSIS**, making the installation process simple and intuitive.

## Features

- **Cross-Platform**: Works on Linux, Windows, and macOS.
- **User Login**: Users can securely log in to the app.
- **File Integration**: Right-click on supported files (as configured in the code) and directly open the FaxSender app with the file pre-attached for sending.
- **Automatic Attachment**: Files are automatically attached to the fax sending window, with pre-filled information.
- **API Integration**: The app interacts with external APIs to manage fax sending.
- **Installer**: The app includes an installer built with **NSIS** for easy installation on Windows.

## Technologies Used

- **Go** (Backend)
- **Fyne** (UI Framework)
- **NSIS** (Installer for Windows)
- **Docker** (For building Docker images)
- **RPM/DEB Packaging** (For Linux deployment)
- **Makefile** (For automating builds and deployment)

## Getting Started

To get started with **FaxSender**, follow the steps below:

### Prerequisites

Ensure you have the following tools installed:

- **Go** (Version 1.21.4 or later)
- **Make**
- **Docker** (Optional for Docker-based builds)
- **NSIS** (For building the Windows installer)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/<your-username>/FaxSender.git
   cd FaxSender
   
2. Build the project:

   ```bash
   make init

3. Run the application (UI for sending faxes):

   ```bash
   make run_ui

4. To build the installer for Windows:

   ```bash
   make deploy_windows_installer

5. To build and deploy for Linux (Debian/Ubuntu or CentOS/RPM-based distros):

*  For Debian/Ubuntu-based systems:

   ```bash
   make deploy_deb

*  For RPM-based systems (CentOS, Fedora):
  
   ```bash
   make deploy_rpm

### File Integration

* Right-click on a supported file type (as configured in the code) and select FaxSender from the context menu.
  
* The file will automatically be attached to the fax window, where you can enter additional information before sending the fax.
  
### Running Tests
You can run the tests with the following command:

    make test

### Docker Support

If you prefer to build the application inside Docker containers, you can use the following targets:

* Build the Docker image for Debian/Ubuntu:
  
  ```bash
  make docker_build_deb

* Build the Docker image for CentOS/RPM:

  ```bash
  make docker_build_rpm

* Run the Docker container for Debian-based systems:
    
  ```bash
  make docker_run_deb

* Run the Docker container for RPM-based systems: 

  ```bash
  make docker_run_rpm

### Deployment
You can deploy the app to various platforms with the following commands:

* Deploy Linux Daemon:

  ```bash
  make deploy_linux_daemon

* Deploy Linux UI:

  ```bash
  make deploy_linux_ui

* Deploy Windows UI:

  ```bash
  make deploy_windows_ui

* Deploy Windows Daemon:

  ```bash
  make deploy_windows_daemon

### License
This project is licensed under the MIT License - see the LICENSE file for details.







    

    






