dataset.nowac = read.table("helper/exprs.csv", sep=",", header=TRUE)

add <- function(a,b) {
  a
  b
  return (a+b)
}

sub <- function(a,b){
  return (a-b)
}
