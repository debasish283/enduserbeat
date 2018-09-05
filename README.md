# Enduserbeat

Welcome to Enduserbeat.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/Manjukb/enduserbeat`

## Getting Started with Enduserbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Enduserbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Enduserbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/Manjukb/enduserbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Enduserbeat run the command below. This will generate a binary
in the same directory with the name enduserbeat.

```
make
```


### Run

To run Enduserbeat with debugging output enabled, run:

```
./enduserbeat -c enduserbeat.yml -e -d "*"
```


### Test

To test Enduserbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Enduserbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Enduserbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/Manjukb/enduserbeat
git clone https://github.com/Manjukb/enduserbeat ${GOPATH}/src/github.com/Manjukb/enduserbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.


## Fields Description (Hardware)

| Fields | Description |
| --- | --- |
| timestamp | elastic timestamp |
| FreeRam | available ram  memory |
| Fstype | file system type |
| HddFree | Hard disc free space in bytes |
| HddTotal |  Hard disc total space in bytes |
| HddUsed |  Hard disc used space in bytes |
| HddUsedPercent |  Hard disc used percentage |
| Hostname |  host name of the machine |
| KernelVersion |  Windows kernel version of OS |
| Model |     CPU model |
| ModelName |               CPU model name |
| OS |                      Operating system |
| Platform |                Operating system Platform |
| Ram |                     RAM of the system |
| RamUsedPercent |          RAM used percent in a particular time |
| Uptime |                  System up time |
| booTime |                 time of beat start --need to check |
| host |                    Network host name of the machine |
| lastPatchUpdateDate |     Last patch update date of operating system |
| manufacturer |            manufacturer of the machine |
| modelname |               model name of the machine |
| servicepack |             service pack version |
| type |                    type of beat (software/Hardware) |
| username |                user under which beat is running |


## Fields Description (Software)

| Fields | Description |
| --- | --- |
| Architecture | software Architecture |
| Hostname |     Network host name of the machine |
| Version |      software version   |
| host |         Network host name of the machine |
| name |         software name |
| tags |         beat tag |
| type |         type of beat (software/Hardware) |
