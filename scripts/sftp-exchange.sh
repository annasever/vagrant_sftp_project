#!/bin/bash

DATE=$(date +'%Y-%m-%d')
TIME=$(date +'%H:%M:%S')
SERVER_NAME=$(hostname)

FILE_CONTENT="${DATE}, ${TIME}, ${SERVER_NAME}"
FILE_NAME="${SERVER_NAME}_file_${DATE}_${TIME}.txt"

# Define SFTP servers
SFTP_SERVERS=("10.0.0.201" "10.0.0.202" "10.0.0.203")
USERNAME="sftp"

# Create a temporary file with the content
echo $FILE_CONTENT > /tmp/$FILE_NAME

CURRENT_IP=$(hostname -I | awk '{print $2}')

# Loop through each SFTP server and upload the file
for SERVER in "${SFTP_SERVERS[@]}"; do
    if [ "$SERVER" != "$CURRENT_IP" ]; then
        sftp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null  $USERNAME@$SERVER <<EOF
        cd uploads
        put /tmp/$FILE_NAME
        bye
EOF
    fi 
done

# Clean up the temporary file
rm /tmp/$FILE_NAME
echo "File created and uploaded to the SFTP servers."
