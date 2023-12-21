#!/bin/bash
#
#SBATCH --mail-user=nishchaykarle@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=project2-time
#SBATCH --output=/home/nishchaykarle/parallel/project-2/project-2-NishchayKarle/proj2/benchmark/time.txt
#SBATCH --error=/home/nishchaykarle/parallel/project-2/project-2-NishchayKarle/proj2/benchmark/%j.%N.stderr
#SBATCH --chdir=/home/nishchaykarle/parallel/project-2/project-2-NishchayKarle/proj2/editor/
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=2000
#SBATCH --exclusive
#SBATCH --time=03:45:00


module load golang/1.16.2
for i in {1..5}
do
    go run editor.go small
done

for i in {1..5}
do
    go run editor.go mixture
done

for i in {1..5}
do
    go run editor.go big
done

for i in {1..5}
do
    go run editor.go small pipeline 2
done

for i in {1..5}
do
    go run editor.go mixture pipeline 2
done

for i in {1..5}
do
    go run editor.go big pipeline 2
done

for i in {1..5}
do
    go run editor.go small pipeline 4
done

for i in {1..5}
do
    go run editor.go mixture pipeline 4
done

for i in {1..5}
do
    go run editor.go big pipeline 4
done

for i in {1..5}
do
    go run editor.go small pipeline 6
done

for i in {1..5}
do
    go run editor.go mixture pipeline 6
done

for i in {1..5}
do
    go run editor.go big pipeline 6
done

for i in {1..5}
do
    go run editor.go small pipeline 8
done

for i in {1..5}
do
    go run editor.go mixture pipeline 8
done

for i in {1..5}
do
    go run editor.go big pipeline 8
done

for i in {1..5}
do
    go run editor.go small pipeline 12
done

for i in {1..5}
do
    go run editor.go mixture pipeline 12
done

for i in {1..5}
do
    go run editor.go big pipeline 12
done

for i in {1..5}
do
    go run editor.go small bsp 2
done

for i in {1..5}
do
    go run editor.go mixture bsp 2
done

for i in {1..5}
do
    go run editor.go big bsp 2
done

for i in {1..5}
do
    go run editor.go small bsp 4
done

for i in {1..5}
do
    go run editor.go mixture bsp 4
done

for i in {1..5}
do
    go run editor.go big bsp 4
done

for i in {1..5}
do
    go run editor.go small bsp 6
done

for i in {1..5}
do
    go run editor.go mixture bsp 6
done

for i in {1..5}
do
    go run editor.go big bsp 6
done

for i in {1..5}
do
    go run editor.go small bsp 8
done

for i in {1..5}
do
    go run editor.go mixture bsp 8
done

for i in {1..5}
do
    go run editor.go big bsp 8
done

for i in {1..5}
do
    go run editor.go small bsp 12
done

for i in {1..5}
do
    go run editor.go mixture bsp 12
done

for i in {1..5}
do
    go run editor.go big bsp 12
done

python3 ../benchmark/generate_graph.py