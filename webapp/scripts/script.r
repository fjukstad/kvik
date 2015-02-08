options(width=10000) 

exprs <- function(genes) { 
    len = length(genes) 
    res = runif(len, 0.0, 100.0)
    return (res)
}

fc <- function(genes) { 
    len = length(genes) 
    res = runif(len, 0.0, 1.0)
    return (res)
}
