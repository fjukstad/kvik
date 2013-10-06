
/* 
 * Ok, the $.get() thing don't need anything except theurl..
 */ 

$(function() {
    $('#pathwayFormSubmit').click(function(e){
        e.preventDefault();
        $.get(getPathways(),function(data) {
              $('.result').html(data);
            });
        }); 
});



function getPathways() {
    console.log("gene url in progress"); 
    var pathways = $('#geneSelect').serialize()
    // var result = 'http://localhost:8080/api/dataset/selectGenes/'+genes;
    var result = "/demo/"+genes
    
    window.open(result,'name','toolbar=0,status=0');
    
    return result; 
}



