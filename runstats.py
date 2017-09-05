#!/usr/bin/python3

import sys
import yaml
import subprocess

def processYamlFile(file):
    with open(file, 'r') as stream:
        try:
            return yaml.load(stream)
        except err:
            print(err)
            return {}

def sumDictElemnts(a, b):
    ret = {};
    for key in a:
        if isinstance(a[key], dict):
            ret[key] = sumDictElemnts(a[key], b[key])
        else:
            ret[key] = a[key]+b[key]
    return ret;

def divideDictByNum(a, n):
    ret = {};
    for key in a:
        if isinstance(a[key], dict):
            ret[key] = divideDictByNum(a[key], n)
        else:
            ret[key] = a[key]/n
    return ret;

if __name__ == "__main__":
    iterations = int(sys.argv[1])
    execname = sys.argv[2]

    statfiles = map(lambda x: str(execname)+'.stats'+str(x)+'.yml', range(iterations))

    stats = []
    for f in statfiles:
        subprocess.call(['./qdi-riscv', '-memfile='+execname, '-statsfile='+f]+sys.argv[:3])
        stats.append(processYamlFile(f));
        
    count = 0
    mean = {}
    for a in stats:
        if count == 0:
            mean = a
        else:
            mean = sumDictElemnts(mean, a);
        count += 1
        
    mean = divideDictByNum(mean, count)
    with open(execname+'.meanstats.yml', 'w') as stream:
        stream.write(yaml.dump(mean))

    print("Bubbles per instruction: "+str((mean['bubbles']/mean['decoded'])*100)+'%');
    print("Branch misprediction ratio: "+str((mean['cancelled']/mean['retired'])*100)+'%')

    total = mean['bubbles']+mean['decoded']
    print("Bypass unit usage: "+str((mean['unit']['bypass']/total)*100)+'%');
    print("Adder unit usage: "+str((mean['unit']['adder']/total)*100)+'%');
    print("Logic unit usage: "+str((mean['unit']['logic']/total)*100)+'%');
    print("Shifter unit usage: "+str((mean['unit']['shifter']/total)*100)+'%');
    print("Branch unit usage: "+str((mean['unit']['branch']/total)*100)+'%');
    print("Memory unit usage: "+str((mean['unit']['memory']/total)*100)+'%');
