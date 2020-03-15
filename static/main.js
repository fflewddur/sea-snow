// This falling snow effect is heavily based on https://codepen.io/loktar00/pen/CHpGo

var refreshInterval = 1000 * 20, // refresh every 20 seconds
    flakes = [],
    canvas = document.getElementById("canvas"),
    ctx = canvas.getContext("2d"),
    flakeCount = 600,
    mouseX = -100,
    mouseY = -100,
    isSnowing = false

canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

$(document).ready(function () {
    initSnow();

    $('body').mousemove(function (e) {
        mouseX = e.pageX;
        mouseY = e.pageY;
    });

    $('body').mouseout(function (e) {
        // If the mouse leaves the page, stop interacting with flakes
        mouseX = -100;
        mouseY = -100;
    });

    $(window).resize(function () {
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
    });

    if (isSnowing) {
        startSnow();
    }
    setTimeout(checkWeather, refreshInterval);
});

function checkWeather() {
    $.get("/api/update", function (data) {
        var pageIsSnowing = $('.snow').length > 0;

        if (data.snowing) {
            if (!pageIsSnowing) {
                $('#status').html('<p class="snow">Woah, it\'s snowing in Seattle!!</p>');
                startSnow();
            }
        } else {
            if (pageIsSnowing) {
                $('#status').html('<p class="nosnow">Nope, it doesn\'t look like it\'s snowing in Seattle.</p>');
                stopSnow();
            }
        }
    }, "json");
    setTimeout(checkWeather, refreshInterval)
}

(function () {
    var requestAnimationFrame = window.requestAnimationFrame || window.mozRequestAnimationFrame || window.webkitRequestAnimationFrame || window.msRequestAnimationFrame ||
        function (callback) {
            window.setTimeout(callback, 1000 / 60);
        };
    window.requestAnimationFrame = requestAnimationFrame;
})();

function startSnow() {
    resetFlakes();
    isSnowing = true;
    $('#canvas').show();
    snow();
}

function stopSnow() {
    isSnowing = false;
}

function snow() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    var flakesAreFalling = false;

    for (var i = 0; i < flakeCount; i++) {
        var flake = flakes[i],
            x = mouseX,
            y = mouseY,
            minDist = 150,
            x2 = flake.x,
            y2 = flake.y;

        // Let flakes finish falling before we stop the animation
        if (!isSnowing && (flake.y >= canvas.height || flake.y <= -10)) {
            continue;
        }
        flakesAreFalling = true;

        var dist = Math.sqrt((x2 - x) * (x2 - x) + (y2 - y) * (y2 - y))

        if (dist < minDist) {
            var force = minDist / (dist * dist),
                xcomp = (x - x2) / dist,
                ycomp = (y - y2) / dist,
                deltaV = force / 2;

            flake.velX -= deltaV * xcomp;
            flake.velY -= deltaV * ycomp;

        } else {
            flake.velX *= .98;
            if (flake.velY <= flake.speed) {
                flake.velY = flake.speed
            }
            flake.velX += Math.cos(flake.step += .05) * flake.stepSize;
        }

        ctx.fillStyle = "rgba(255,255,255," + flake.opacity + ")";
        flake.y += flake.velY;
        flake.x += flake.velX;

        if (flake.y >= canvas.height) {
            reset(flake);
        }

        if (flake.x >= canvas.width || flake.x <= 0) {
            reset(flake);
        }

        ctx.beginPath();
        ctx.arc(flake.x, flake.y, flake.size, 0, Math.PI * 2);
        ctx.fill();
    }
    if (flakesAreFalling) {
        requestAnimationFrame(snow);
    } else if (!isSnowing) {
        $('#canvas').hide();
    }
};

function resetFlakes() {
    for (var i = 0; i < flakeCount; i++) {
        reset(flakes[i])
    }
}

function reset(flake) {
    flake.x = Math.floor(Math.random() * canvas.width);
    flake.y = -1 * Math.floor(Math.random() * canvas.height);
    flake.size = (Math.random() * 3) + 2;
    flake.speed = (Math.random() * 1) + 0.75;
    flake.velY = flake.speed;
    flake.velX = 0;
    flake.opacity = (Math.random() * 0.5) + 0.3;
}

function initSnow() {
    for (var i = 0; i < flakeCount; i++) {
        flake = {
            x: 0,
            y: 0,
            opacity: 0,
            size: 0,
            speed: 0,
            stepSize: (Math.random()) / 30,
            step: 0,
            velX: 0,
            velY: 0
        }
        reset(flake)
        flakes.push(flake);
    }
    isSnowing = $('.snow').length > 0;
};
