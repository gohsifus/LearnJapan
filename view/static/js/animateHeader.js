(function(){
    setInterval(function(){
        let header = document.getElementById('header')
        header.style.backgroundImage = "url('/static/images/" + getRandomInt(1, 4) + ".jpg')";
    }, 6000)
})();

function getRandomInt(min, max) {
        return Math.floor(Math.random() * (max - min + 1)) + min;
}
