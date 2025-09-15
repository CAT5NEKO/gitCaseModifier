# gitCaseModifier

Normally, Git and Github does not detect filename changes that differ only by case 🙀🙀🙀🙀🙀  
With this tool, you can easily apply such updates (e.g., `file.txt` → `File.txt`) in bulk via a simple list file.


## Usage

First, create a simple text/markdown file that lists the files you want to rename.   

Each line should contain the **current filename** followed by the **new filename**, separated by a space.

Example:

```txt
# Rename rules
fuga.ts Fuga.ts
hoge/fuga.vue hoge/Fuga.vue
docs/readme.md docs/README.md
```
Then, build the tool:    

```bash
go build
```

Once built, run it by specifying the repository path and the rename list file:  


```bash
./gitCaseModifier <repo_path> <csv_file>
```