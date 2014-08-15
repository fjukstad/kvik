
/* 
 * Ok, the $.get() thing don't need anything except theurl..
 */ 

$(function() {
    $('#geneFormSubmit').click(function(e){
        e.preventDefault();
        $.get(getGenes(),function(data) {
              $('.result').html(data);
            });
        }); 
});



function getGenes() {
    var genes = $('#geneSelect').serialize()
    // var result = 'http://localhost:8080/api/dataset/selectGenes/'+genes;
    var result = "/demo/"+genes
    
    window.open(result,'name','toolbar=0,status=0');
    
    return result; 
}


