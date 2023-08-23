function data() {
  return {
    displayMode: 'day',
    toggleDisplayMode: function () {

    },
    isPlaying: false,
    time: 0,
    timeStr: '00:00:00',
    duration: 1,
    durationStr: '00:00:00',
    initPlayer: function () {
      this.duration = $refs.audio.duration
      this.durationStr = this.sec2time(this.duration)
      this.updateProgress()
    },
    togglePlayer: function () {
      if (this.isPlaying) {
        this.isPlaying = false
        $refs.audio.pause()
      } else {
        this.isPlaying = true
        $refs.audio.play()
      }
    },
    sec2time: function (sec) {
      sec = Math.floor(sec)
      var hours   = Math.floor(sec / 3600)
      var minutes = Math.floor((sec - (hours * 3600)) / 60)
      var seconds = sec - (hours * 3600) - (minutes * 60)

      if (hours   < 10) {hours   = '0'+hours;}
      if (minutes < 10) {minutes = '0'+minutes;}
      if (seconds < 10) {seconds = '0'+seconds;}
      
      return hours+':'+minutes+':'+seconds
    },
    updateProgress: function () {
      this.time = $refs.audio.currentTime
      this.timeStr = this.sec2time(this.time)
      // console.log('Time: ' + this.time)
      // console.log('Duration: ' + this.duration)
      // console.log('Old value: ' + $refs.progress.value)
      //$refs.progress.value = (this.time / this.duration) * 1000
      //$refs.progress.style.width = 10+Math.round((this.time / this.duration) * $refs.progressbar.offsetWidth)
      $refs.progress.style.backgroundSize = (this.time / this.duration) * 100 + '% 100%'
    },
    seek: function () {
      newPlace = $refs.progress.value
      newTime = ( this.duration / 1000 ) * newPlace
      $refs.audio.currentTime = newTime
    },
    input: function () {
      newPlace = $refs.progress.value
      console.log('New place: ' + newPlace)
      console.log('Style: ' + $refs.progress.style.backgroundSize)
    }
  };
}

function init() {

}

Date.prototype.getLocalTime = function() {
  h = (h = this.getHours()) < 10 ? ('0' + h) : h;
  m = (m = this.getMinutes()) < 10 ? ('0' + m) : m;
  s = (s = this.getSeconds()) < 10 ? ('0' + s) : s;
  return h + ":" + m + ":" + s;
}

function sec2time(sec) {
  sign = (sec < 0)?'-':''
  sec = Math.abs(Math.floor(sec))
  var hours   = Math.floor(sec / 3600)
  var minutes = Math.floor((sec - (hours * 3600)) / 60)
  var seconds = sec - (hours * 3600) - (minutes * 60)

  if (hours   < 10) {hours   = '0'+hours;}
  if (minutes < 10) {minutes = '0'+minutes;}
  if (seconds < 10) {seconds = '0'+seconds;}

  return sign+hours+':'+minutes+':'+seconds
}