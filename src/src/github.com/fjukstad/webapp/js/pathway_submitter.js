

$(function() {
    $('#pathwayFormSubmit').click(function(e){
        e.preventDefault();
        $.get(getPathways(),function(data) {
              $('.result').html(data);
            });
        }); 
});



function getPathways() {
    var pathways = $('#pathwaySelect').serialize()
    var result = "/demo/"+pathways
    window.location.assign(result)
    
   // window.open(result,'name','toolbar=0,status=0,height=700,width=1300');
    
    return result; 
}

function ReadableInput() {
    var sel = document.getElementById('pathwaySelect');
    var opts = sel.options;
    for(var opt, j = 0; opt = opts[j]; j++) {
        opt.innerHTML = GetPathwayName(opt.innerHTML)
    }
} 


function GetPathwayName(id) {
    var baseURL = "http://"+window.location.hostname+":8080/info/pathway/"
        url =  baseURL+id+"/name"

    var name
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            name = data;
        }
    });
    return name;

}

// When window load fetch names for the pathways in the list.
window.onload = function() {
    ReadableInput() 
}



