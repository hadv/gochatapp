'use strict';

angular.module('chatWebApp')
  .factory('socket', ['socketFactory', function(socketFactory) {
    var socket = socketFactory({
      prefix: 'mana:'
    });
    socket.forward('error');
    return socket;
  }]);
