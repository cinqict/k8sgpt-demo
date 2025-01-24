
# Demo plan

- configure
- local llm reply
- cloud llm reply

- extra:
  - ?

- leuke bug
  - bug die pas op te lossen



# Demo setup

## Setup cluster

```shell

kind create cluster

kubectl apply -f ./k8s-bugs/


```

## Ollama on Ubuntu

Recommended: enable usage of gpu inside ollama container

```shell
# Setup Ollama
ollama pull llama3.1

docker run -it --rm \
  --device /dev/dri \
  -e DISPLAY=$DISPLAY \
  -v /tmp/.X11-unix:/tmp/.X11-unix \
  -v ollama:/root/.ollama \
  -e OLLAMA_DEBUG=1 \
  -p 11434:11434 \
  --name ollama \
  ollama/ollama

docker exec -it ollama ollama run llama3.1

# Configure k8sgpt with ollama backend
k8sgpt auth add --backend ollama --baseurl http://localhost:11434 --model llama3.1

# Make it the default backend
k8sgpt auth default --provider ollama

k8sgpt analyze --explain

```

### Troubleshooting

```shell
cat ~/.config/k8sgpt/k8sgpt.yaml

```


## Use your gpu

### Mesa Intel Graphics

On Ubuntu

```shell
# Install Mesa / gpu tools
sudo apt install mesa-utils intel-opencl-icd
sudo apt install intel-gpu-tools intel-media-va-driver-non-free

# Allow Access to X Server
xhost +local:

# Verify
docker run -it --rm \
    --device /dev/dri \
    -e DISPLAY=$DISPLAY \
    -v /tmp/.X11-unix:/tmp/.X11-unix \
    ubuntu bash

ls /dev/dri
apt update
apt install mesa-utils
glxinfo | grep "OpenGL renderer"

```



