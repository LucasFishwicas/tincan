# Pull Node image
FROM node:21

# Install wscat
RUN npm install -g wscat

# Copy local bash script
COPY ./tincan-bash.sh .

# Convert to an executable
RUN chmod +x tincan-bash.sh

# Run executable
CMD ["./tincan-bash.sh"]
