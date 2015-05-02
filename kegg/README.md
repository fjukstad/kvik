# KEGG package
This is the KEGG package Kvik uses to fetch information about pathways, genes
and compounds. It ises the KEGG HTTP REST API, so it will probably take some
time. However it caches all requests to the a cache/ folder, so the second time
it fetches something from the REST API it will just use the contents from the
cache/ folder. Be aware that the cache is never flushed. NEVER. 

# Tools
There are some tools in the tools/ directory. Go have fun. 


