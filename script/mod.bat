@echo off

cd ../

@rem ALib
cd ./src/ALib
go mod tidy
cd ../../

@rem DockerCTL
cd ./src/DockerCTL
go mod tidy
cd ../../

@rem DockerCTL
cd ./src/DomainHealthCTL
go mod tidy
cd ../../

@rem GitlabCTL
cd ./src/GitlabCTL
go mod tidy
cd ../../

@rem KubernetesCTL
cd ./src/KubernetesCTL
go mod tidy
cd ../../

@rem NacosCTL
cd ./src/NacosCTL
go mod tidy
cd ../../