window.onload = function() { 
    var cache = {};
    $( "#input" ).autocomplete({
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
        });
      }
    });
}; 
