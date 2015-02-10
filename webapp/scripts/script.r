options(width=10000) 
rnorm2 <- function(n,mean,sd) {
    mean+sd*scale(rnorm(n))
}

exprs <- function(gene) { 
    #len = length(genes) 
    len = 50
    res <- rnorm2(50, 0,0.45) 
    return (c(res))
}

fc <- function(genes) { 
    len = length(genes) 
    res = runif(len, -1.0, 1.0)
    return (res)
}
