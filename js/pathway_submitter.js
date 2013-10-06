

$(function() {
    $('#pathwayFormSubmit').click(function(e){
        e.preventDefault();
        $.get(getPathways(),function(data) {
              $('.result').html(data);
            });
        }); 
});



function getPathways() {
    console.log("Fetching those pathways for you"); 
    var pathways = $('#pathwaySelect').serialize()
    var result = "/demo/"+pathways
    
    window.open(result,'name','toolbar=0,status=0');
    
    return result; 
}



