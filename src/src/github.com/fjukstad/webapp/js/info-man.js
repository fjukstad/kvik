// js for retrieving info about specific pathway/gene/compund etc.

function GetInfo(name) {


    var a = name;
    var b = a.replace(/ /g, '+').toLowerCase();

    var baseURL = "http://"+window.location.hostname+":8080/info";
    var infoType = "all";

    url = baseURL + "/" + b + "/" + infoType;

    var info; 
    
    $.ajax({
        async: false,
        cache: true,
        type: "GET",
        url: url,
        dataType: "text",
        success: function(data) {
            info = data;
        }
    });

    return jQuery.parseJSON(info);

}
