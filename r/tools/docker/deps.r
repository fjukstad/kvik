packages <- c("ggplot2", "Hmisc", "Rcpp", "roxygen2", "jsonlite", "igraph", "dplyr", "multicore", "colorspace") 
install.packages(packages, repos='http://cran.rstudio.com/')

source("http://bioconductor.org/biocLite.R")
pkgs <- c("Biobase" ,"DBI" ,"RSQLite" ,"AnnotationDbi" ,"GO.db" ,"RColorBrewer" ,"latticeExtra" ,"colorspace" ,"munsell" ,"scales" ,"ggplot2" ,"Hmisc" ,"reshape" ,"preprocessCore", "impute", "WGCNA" ,"illuminaHumanv3.db" ,"illuminaHumanv4.db" ,"animation" ,"limma")
biocLite(pkgs, ask=FALSE) 
