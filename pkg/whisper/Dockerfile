FROM continuumio/miniconda3

WORKDIR /app

# Create the environment:
COPY env.yml .
RUN conda env create -f env.yml

# Make RUN commands use the new environment:
RUN echo "conda activate myenv" >> ~/.bashrc
SHELL ["/bin/bash", "--login", "-c"]

# Make RUN commands use the new environment:
SHELL ["conda", "run", "-n", "myenv", "/bin/bash", "-c"]

# The code to run when container is started:
COPY whisper.py .
ENTRYPOINT ["python", "whisper.py"]