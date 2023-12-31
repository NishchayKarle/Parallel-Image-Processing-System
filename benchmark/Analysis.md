## SPEEDUP GRAPHS

- **PIPELINE**
  - ![](/proj2/benchmark/speedup-pipeline.png)

- **BSP**
  - ![](/proj2/benchmark/speedup-bsp.png)

---

## A brief description of the project.

In this project, a set of images needs to be converted into resultant images by applying a list of effects. The effects are stored in a text file, and each image needs to go through each effect in the file.

### Sequential Version:

For every size (small/mixture/large), the effects.txt file is read. Based on the JSON, the image is loaded, and the list of effects is applied. Once the effects are applied, the resultant image is saved in an output path. After all the actions on an image are done, the next image is picked from effects.txt, and the process is repeated.


### Pipeline Version:

In the pipeline version, the task creator puts all the images into a channel. One worker thread works on one image, and it reads the image from the channel. 

After that, the worker thread spawns multiple mini-workers that apply all effects on the bounds assigned to that mini-worker. Here, channels are used to implement wait groups.

Once all the threads are done applying an effect, the image is swapped to be ready for the next effect. This step is important to ensure that all the threads apply the effect before moving on to the next one. 

The filtered image is fed into a channel, and all the images are saved in their output path. Here, wait groups are implemented using channels.

### BSP Version:

In the BSP version, a barrier `barrier.go` is implemented.

The `Run BSP Worker` implementation works differently based on the thread ID. If the ID is the master ID (in our case, ID = threadCount-1), the effects.txt file is read, and all the worker threads including the master will work on one chunk of this image. 

In the BSP version, we also wait until all worker threads have applied one effect on their assigned bounds (bounds are assigned based on the ID) before moving on to the next effect. 

Once all the effects on an image are applied, the image is pushed to an array. The master routine (threadCount-1) will then spawn new threads to save the images. In order to wait for threads, barriers are used.

## Instructions on how to run your testing script.

---

The sbatch script is located in `proj2/benchmark` named `benchmark-proj2.sh`.

The generate_graph.py requires the correct path of benchmark/time.txt (will be generated by the sbatch job) to be provided on line 14

To run the sbatch script, follow these steps:

```
cd proj2/benchmark
sbatch benchmark-proj2.sh
```

The results will be generated in data.txt file inside the `benchmark` folder.

## Explanation of graph results

- What are the hotspots and bottlenecks in your sequential program?

  - The following are the bottlenecks and hotspots in the sequential program:

    - Although the images in the list `effects.txt` are not dependent on each other for actions, we are waiting for one image to finish to start another.
    - Similarly, the rows in the images are not dependent on one another to apply effects, but still, while doing sequentially, we are waiting for one row to finish before moving on to another.
    - While the main go routine is saving, it's waiting for the image to be saved in the file before it can start an operation on another image.

- Which parallel implementation is performing better? Why do you think it is?

    - The BSP with 12 threads is performing the best among all three sequential, pipeline, and BSP. Although the speedup was even higher if I used the maximum possible threads, which is 128 on the Linux cluster.

In both pipeline and BSP, the image is divided into chunks, and each Go routine works on one chunk
