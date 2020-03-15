var refreshInterval = 1000 * 20; // refresh every 20 seconds

$(document).ready(function () {
    setTimeout(checkWeather, refreshInterval);
});

function checkWeather() {
    $.get("/api/update", function(data) {
        var pageIsSnowing = $('.snow').length > 0;

        if (data.snowing) {
            if (!pageIsSnowing) {
                $('#status').html('<p class="snow">Woah, it\'s snowing in Seattle!!</p>');
                // TODO start snow animation
            }
        } else {
            if (pageIsSnowing) {
                $('#status').html('<p class="nosnow">Nope, it doesn\'t look like it\'s snowing in Seattle.</p>');
                // TODO stop snow animation
            }
        }
    }, "json");
    setTimeout(checkWeather, refreshInterval)
}