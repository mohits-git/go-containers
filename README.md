# Go Containers

A simple container implementation in Go.

## What is a container?

A container is a way to isolate a process from the rest of the system. It is a way to run a process in a sandboxed environment.

## How does a container work?

Containers are made possible by the 3 major Linux kernel features: 
- Namespaces: Namespaces allow a process to have its own view of the system. This means that a process can have its own view of the filesystem, network, and other system resources.
- Chroot: Chroot allows a process to change its root directory. This means that a process can have its own view of the filesystem.
- Cgroups: Cgroups allow a process to have resource limits. This means that a process can have limits on how much CPU, memory, and other resources it can use.

## What is this project's goal?

This project aims to implement a simple container in Go. We will create a ubuntu container that runs a simple shell.

## Project Setup

### Requirements for project

- Go 1.16 or higher
- Linux
- Root access
- ubuntu base image

> I have ubuntu cgroup v2 enabled. If you have cgroup v1 enabled, you can change the cgroup path (`/sys/fs/cgroup/pids`) in the code.

### Steps to run the project

1. Clone the repository
```bash
git clone https://github.com/mohits-git/go-containers.git
cd go-containers
```

2. Run the following command to build the container:
```bash
go build -o go-container main.go
```

3. Extract the ubuntu base image to ./rootfs:
```bash
mkdir rootfs
sudo tar -xzf ubuntu-base-24.10-base-arm64.tar.gz -C rootfs
```
> You can download the ubuntu base image from [here](https://cdimage.ubuntu.com/ubuntu-base/releases/24.10/release/) according to your architecture.

4. Run the following command to run the container:
```bash
sudo ./go-container run /bin/bash
```

Now try running the fork bomb in the container >_> :
```bash
:(){ :|:& };:
```
It should not affect the host system. Your container will have max 20 processes limit :D.

## References

[Containers from Scratch - Liz Rice - GOTO 2018](https://youtu.be/8fi7uSYlOdc?si=f_Bl-sKthnjykthk)
