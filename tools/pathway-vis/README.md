# KEGG Pathway visualization library
JavaScript library to visualize KEGG pathways using [d3](d3js.org). 

# Usage 
First run the server or host your own that can serve JSON of the KEGG graphs. 

```
go run main.go
```

then visit `localhost:8080` to view the pathway. In HTML: 

```
<html>
    <div id="pathway-container"></div> 
    <script src="d3.v3.min.js"></script> 
    <script src="pathway.js"></script>
    <script>
        pathway("hsa05200", "#pathway-container", 1000, 1000);
    </script>
</html>
```
