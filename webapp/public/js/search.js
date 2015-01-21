window.onload = function() { 
    var cache = {};
    $( "#input" ).autocomplete({
      minLength: 2,
      source: function( request, response ) {
        var term = request["term"];
        console.log("request:", request);
        console.log("term", term) 
        if ( term in cache ) {
          response( cache[ term ] );
          return;
        }
        $.getJSON("/search", request, function( data, status, xhr ) {
            console.log("response:", data) 
            console.log("cache",cache) 
            cache[ term ] = data;
            response( data );
        });
      }
    });
}; 
