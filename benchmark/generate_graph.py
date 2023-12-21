import matplotlib.pyplot as plt
import numpy as np


def divide(arr):
    val = arr[0]
    for i in range(len(arr)):
        arr[i] /= val
        arr[i] = 1 / arr[i]
        arr[i] = round(arr[i], 2)


def generateGraphs():
    f = open("time.txt", "r")

    seq = [0.0, 0.0, 0.0]
    for i in range(3):
        time = 0
        for j in range(5):
            time += float(f.readline().strip("\n"))
        seq[i] = time / 5

    pipe_small = [seq[0]]
    pipe_mixture = [seq[1]]
    pipe_big = [seq[2]]

    for t in {2, 4, 6, 8, 12}:
        for j in range(3):
            time = 0
            for i in range(5):
                time += float(f.readline().strip("\n"))

            if j == 0:
                pipe_small.append(time / 5)

            elif j == 1:
                pipe_mixture.append(time / 5)

            else:
                pipe_big.append(time / 5)

    bsp_small = [seq[0]]
    bsp_mixture = [seq[1]]
    bsp_big = [seq[2]]

    for t in {2, 4, 6, 8, 12}:
        for j in range(3):
            time = 0
            for i in range(5):
                time += float(f.readline().strip("\n"))

            if j == 0:
                bsp_small.append(time / 5)

            elif j == 1:
                bsp_mixture.append(time / 5)

            else:
                bsp_big.append(time / 5)

    divide(pipe_small)
    divide(pipe_mixture)
    divide(pipe_big)

    divide(bsp_small)
    divide(bsp_mixture)
    divide(bsp_big)

    xpoints = np.array([1, 2, 4, 6, 8, 12])
    plt.title("SPEEDUP GRAPH (PIPELINE)")

    ypoints = np.array(pipe_small)
    plt.plot(xpoints, ypoints, marker="o", label="SMALL")

    ypoints = np.array(pipe_mixture)
    plt.plot(xpoints, ypoints, marker="o", label="MIXTURE")

    ypoints = np.array(pipe_big)
    plt.plot(xpoints, ypoints, marker="o", label="BIG")

    plt.legend()

    plt.xlabel("NUM OF THREADS")
    plt.ylabel("SPEEDUP")
    plt.savefig("speedup-pipeline.png")

    fig = plt.figure()
    plt.figure().clear()
    plt.close()
    plt.cla()
    plt.clf()

    xpoints = np.array([1, 2, 4, 6, 8, 12])
    plt.title("SPEEDUP GRAPH (BSP)")

    ypoints = np.array(bsp_small)
    plt.plot(xpoints, ypoints, marker="o", label="SMALL")

    ypoints = np.array(bsp_mixture)
    plt.plot(xpoints, ypoints, marker="o", label="MIXTURE")

    ypoints = np.array(bsp_big)
    plt.plot(xpoints, ypoints, marker="o", label="BIG")

    plt.legend()

    plt.xlabel("NUM OF THREADS")
    plt.ylabel("SPEEDUP")
    plt.savefig("speedup-bsp.png")


if __name__ == "__main__":
    generateGraphs()
