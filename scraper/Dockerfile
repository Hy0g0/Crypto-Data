# Use an official Python runtime as a base image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the requirements.txt file into the container
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of your scraper code into the container
COPY jina_scraper.py .

# Set the environment variable for unbuffered output
ENV PYTHONUNBUFFERED=1

VOLUME [ "/user-secret" ]
# Command to run your scraper
CMD ["python", "jina_scraper.py"]
