PBFT Simulator
------

This repository contains the Golang code of simple pbft consensus implementation.

- **Original Codes Authorship** 
github.com/tn606024/simplePBFT

- **Modified by** 
@lcy1317 & @leyuwei 

- **Modified for simulations in paper** <br>
Age-of-Information Analysis for Blockchain-Based Mobile Edge Computing 

How to run
------

## Quick Start

1. Run the following command in CMD
```shell script
./simplePBFT
```

2. Wait a few minutes, and the simulation results will be saved in a file containing `res.txt`.

3. To analyze and visualize the PDF of the recorded PBFT latencies, run `analyze_7nodes.py` with Python.

4. If you want to add or remove nodes, you need to modify both `data.go` and `main.go`. Currently, all nodes are created manually.

**It is strongly recommended to delete all `.txt` files and the `Keys` folder after making any changes to the source code.**

**You may encounter node offline errors. If this happens, restart `simplePBFT.exe` to allow the program to continue logging data into the `res.txt` file.**

## Build Your Own Executable

```shell script
go build 
```

### Start Batch Sim

This will automatically start 7 PBFT nodes and a client. 

```shell script
./simplePBFT
```

