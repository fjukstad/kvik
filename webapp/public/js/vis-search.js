function initSearch() { 
    var idmap = {} 
    var cache = {};

    $( "input#search" ).autocomplete({
      minLength: 2,
      source: function( request, response ) {
        var term = request["term"];
        if ( term in cache ) {
          response( cache[ term ] );
          return;
        }
        $.getJSON("/search", request, function( data, status, xhr ) {
            cache[ term ] = data.Terms;
            response( data.Terms );

            if(data.Terms.length > 0){ 
                for(i=0;i<data.Terms.length;i++){
                    pw = data.Terms[i]
                    idmap[ pw.label ] = pw.id
                }
            }
        });
      }
    });

    $("input#search").bind("enterKey", function(e){
        searchterm = $("input#search").val()
        id = idmap[searchterm]
        if(typeof id === 'undefined'){
            swal({
                title: "Could not find what you we're searching for, sorry!",
                text: "You searched for '"+searchterm+"'",
                type: "warning"
            }) 
            return
        } else {  
            pathway(id,"content",0,0) 
        }
    });

    $("input#search").keyup(function(e){
        if(e.keyCode == 13) {
            $(this).trigger("enterKey");
        }
    }) 
}; 

