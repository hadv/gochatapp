'use strict';

angular.module('chatWebApp')
  .factory('socket', ['socketFactory', function(socketFactory) {
    var socket = socketFactory({
      prefix: 'mana:',
      ioSocket: io.connect()
    });
    socket.forward('error');
    return socket;
  }]);
