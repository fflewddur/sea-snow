var refreshInterval = 1000 * 20; // refresh every 20 seconds

$(document).ready(function () {
    console.log("ready");
    setTimeout(checkWeather, refreshInterval);
});

function checkWeather() {
    console.log("checkWeather()")
    // TODO fetch weather conditions via AJAX call
    setTimeout(checkWeather, refreshInterval)
}