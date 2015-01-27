window.onload = function() { 
    var pathwayList = new List('pathway-search', {
        valueNames: ['name'],
        indexAsync: true
    }); 
    var geneList = new List('gene-search', {
        valueNames: ['name'],
        indexAsync: true
    }); 
    $('.search').on("keyup", function(){
        pathwayList.search($(".search").val());
        geneList.search($(".search").val());
    }); 

}; 
