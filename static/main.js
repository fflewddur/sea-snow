var refreshInterval = 1000 * 20, // refresh every 20 seconds
    flakes = [],
    canvas = document.getElementById("canvas"),
    ctx = canvas.getContext("2d"),
    flakeCount = 400,
    mX = -100,
    mY = -100,
    isSnowing = false

canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

$(document).ready(function () {
    initSnow();
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
    isSnowing = true;
    $('#canvas').show();
    snow();
}

function stopSnow() {
    isSnowing = false;
    $('#canvas').hide();
}

function snow() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    for (var i = 0; i < flakeCount; i++) {
        var flake = flakes[i],
            x = mX,
            y = mY,
            minDist = 150,
            x2 = flake.x,
            y2 = flake.y;

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

        if (flake.y >= canvas.height || flake.y <= 0) {
            reset(flake);
        }


        if (flake.x >= canvas.width || flake.x <= 0) {
            reset(flake);
        }

        ctx.beginPath();
        ctx.arc(flake.x, flake.y, flake.size, 0, Math.PI * 2);
        ctx.fill();
    }
    if (isSnowing) {
        requestAnimationFrame(snow);
    }
};

function reset(flake) {
    flake.x = Math.floor(Math.random() * canvas.width);
    flake.y = 0;
    flake.size = (Math.random() * 3) + 2;
    flake.speed = (Math.random() * 1) + 0.5;
    flake.velY = flake.speed;
    flake.velX = 0;
    flake.opacity = (Math.random() * 0.5) + 0.3;
}

function initSnow() {
    for (var i = 0; i < flakeCount; i++) {
        var x = Math.floor(Math.random() * canvas.width),
            y = Math.floor(Math.random() * canvas.height),
            size = (Math.random() * 3) + 2,
            speed = (Math.random() * 1) + 0.5,
            opacity = (Math.random() * 0.5) + 0.3;

        flakes.push({
            speed: speed,
            velY: speed,
            velX: 0,
            x: x,
            y: y,
            size: size,
            stepSize: (Math.random()) / 30,
            step: 0,
            opacity: opacity
        });
    }
    isSnowing = $('.snow').length > 0;
};

canvas.addEventListener("mousemove", function (e) {
    mX = e.clientX,
        mY = e.clientY
});

window.addEventListener("resize", function () {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
})
