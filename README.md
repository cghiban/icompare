#icompare

icompare is a tool for comparing entries between local directory and remote iRods directory

It only looks at the file names, not size or checkusms. If there are differences, you get an output like this:

```
Oh no!! Here's the diff:
   []main.entry{
         ... // 9 identical elements
         {Name: "W05_filtered_L001_R1_001.fastq.gz"},
         {Name: "W06_filtered_L001_R1_001.fastq.gz"},
 -       {Name: "ZZ-test", IsDir: true},
         {Name: "mappingfile.tsv"},
   }
```

Usage:
```bash
$ icompare -h
Usage of ./icompare:
  -l string
        Local path
  -r string
        Remote path
```

Example:
```bash
$ icompare -l ./test -r /remote/irods/path/test
```

##requirements
It needs is from `irods-icommands`.

##todo

- [ ] make it work recursively

