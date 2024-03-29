# STLSDataParser
STLSDataParser is a tool designed to monitor, analyze, and transmit transport layer security analysis results generated by the Secure Transport Layer Scanner (STLS). This service automates the process of capturing JSON files resulting from STLS scans, parses this data, and prepares it to be sent to a messaging or storage system, facilitating integration with other security monitoring tools and systems.

# Features
+ Directory Monitoring: Continuously watches a specified folder for new JSON files generated by STLS.
+ Data Parse: Parses the JSON files, extracts relevant information, and structures the data for easy transmission.
+ Data Sending: Flexible configuration for sending the parsed data to messaging or storage systems, adaptable to the user's needs.

# Configuration
Before starting STLSDataParser, configure the path of the directory to be monitored and the details of the messaging system in the config.yaml file:

# Contributions
Contributions are welcome! If you have improvements or fixes, please fork the repository, make your changes, and send a pull request.

# License
This project is licensed under the MIT License - see the LICENSE file for details.

# Contact us
If you have any questions or suggestions, please open an issue on GitHub.
