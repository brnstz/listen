$(document).ready(function() {
    $.getJSON("api/shows.json", function(data) {
        var shows = $("div#shows");

        for (var i = 0; i < data.length; i++) {
            var newRow = $("<div>", {class: "row"});
           
            var colDate = $("<div>", {class: "col-md-4"});
            $(colDate).text(data[i].starts_formatted);

            var colVenue = $("<div>", {class: "col-md-4"});
            $(colVenue).text(data[i].venue.name);

            var colBand = $("<div>", {class: "col-md-4"});
            $(colBand).text(data[i].bands[0].name);


            shows.append(newRow);
            newRow.append(colDate, colVenue, colBand);
        }
    });
});
