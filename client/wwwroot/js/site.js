// Please see documentation at https://docs.microsoft.com/aspnet/core/client-side/bundling-and-minification
// for details on configuring this project to bundle and minify static web assets.

// Write your JavaScript code.

function ready(fn) {
    if (document.attachEvent ? document.readyState === "complete" : document.readyState !== "loading") {
        fn();
    } else {
        document.addEventListener('DOMContentLoaded', fn);
    }
}

var maxX = 20;
var maxY = 20;
var cycleTimeout = 1;
var boxCache = {};

function onready() {
    let boxes = document.getElementById("boxes");
    let x = 0;
    let y = 0;

    console.log('Starting');

    for (x = 0; x < maxX; x++) {
        let row = document.createElement("div");
        boxes.append(row);
        for (y = 0; y < maxY; y++) {
            let box = document.createElement("div");
            box.classList.add('box');
            box.id = "_x" + x + "_y" + y;
            box.innerText = box.id;
            row.append(box);

            boxCache[box.id] = box;
        }
    }

    let interval = setInterval(cycler, cycleTimeout);
}

var cycleX = 0;
var cycleY = 0;
function cycler() {
    if (cycleX >= maxX) {
        cycleX = 0;
        cycleY++;
    }

    if (cycleY >= maxY) {
        cycleY = 0;
        cycleX = 0;
    }

    let boxId = "_x" + cycleX + "_y" + cycleY;
    let box = boxCache[boxId];

    query(box, cycleX, cycleY);

    cycleX++;
}

function query(box, x, y) {
    var request = new XMLHttpRequest();
    box.classList.add('loading');
    //request.open('POST', '/api/count', true);
    request.open('GET', 'http://localhost:10000', true);
    request.setRequestHeader('Content-Type', 'application/json');

    request.onload = function () {
        if (request.status >= 200 && request.status < 400) {
            // Success!
            var data = JSON.parse(request.responseText);
            box.innerText = data.id;
            box.classList.remove('loading');
        } else {
            console.log('oops ' + request.statusText);

        }
    };

    request.onerror = function () {
        // There was a connection error of some sort
    };

    request.send(JSON.stringify({ x: x, y: y }));
}

ready(onready);
