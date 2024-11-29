FaxSender - Cross-Platform Fax Sending Application

FaxSender is a cross-platform desktop application built with Go for the backend and the Fyne package for the user interface. This application allows users to send faxes directly from their desktop. It supports integration with an installer built using NSIS, making the installation process simple and intuitive.
Features

    Cross-Platform: Works on Linux, Windows, and macOS.
    User Login: Users can securely log in to the app.
    File Integration: Right-click on supported files (as configured in the code) and directly open the FaxSender app with the file pre-attached for sending.
    Automatic Attachment: Files are automatically attached to the fax sending window, with pre-filled information.
    API Integration: The app interacts with external APIs to manage fax sending.
    Installer: The app includes an installer built with NSIS for easy installation on Windows.

Technologies Used

    Go (Backend)
    Fyne (UI Framework)
    NSIS (Installer for Windows)
    Docker (For building Docker images)
    RPM/DEB Packaging (For Linux deployment)
    Makefile (For automating builds and deployment)

Getting Started

To get started with FaxSender, follow the steps below:
Prerequisites

Ensure you have the following tools installed:

    Go (Version 1.21.4 or later)
    Make
    Docker (Optional for Docker-based builds)
    NSIS (For building the Windows installer)

Installation

    Clone the repository:

git clone https://github.com/<your-username>/FaxSender.git
cd FaxSender

Build the project:

make init

Run the application (UI for sending faxes):

make run_ui

To build the installer for Windows:

make deploy_windows_installer

This will build the Windows installer using NSIS.

To build and deploy for Linux (Debian/Ubuntu or CentOS/RPM-based distros):

    make deploy_deb    # For Debian/Ubuntu-based systems
    make deploy_rpm    # For RPM-based systems (CentOS, Fedora)

File Integration

    Right-click on a supported file type (as configured in the code) and select FaxSender from the context menu.
    The file will automatically be attached to the fax window, where you can enter additional information before sending the fax.

Running Tests

You can run the tests with the following command:

make test

Docker Support

If you prefer to build the application inside Docker containers, you can use the following targets:

    Build the Docker image for Debian/Ubuntu:

make docker_build_deb

Build the Docker image for CentOS/RPM:

make docker_build_rpm

Run the Docker container:

    make docker_run_deb    # For Debian-based systems
    make docker_run_rpm    # For RPM-based systems

Deployment

You can deploy the app to various platforms with the following commands:

    Deploy Linux Daemon:

make deploy_linux_daemon

Deploy Linux UI:

make deploy_linux_ui

Deploy Windows UI:

make deploy_windows_ui

Deploy Windows Daemon:

    make deploy_windows_daemon

License

This project is licensed under the MIT License - see the LICENSE file for details.
