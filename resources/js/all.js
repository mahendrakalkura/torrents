var url = document.querySelector('body').getAttribute('data-url');

var root = document.querySelector('#root');

var oninit = function() {
    var socket = new WebSocket('ws://' + url + '/websockets/');

    socket.addEventListener('open', function(event) {
        socket.send('start');
    });

    socket.addEventListener('close', function(event) {
    });

    socket.addEventListener('error', function(event) {
    });

    socket.addEventListener('message', function(event) {
        component.progress = event.data;
        m.redraw();
        if (event.data === 100) {
            socket.close();
        }
    });
};

var component = {
    progress: 0,
    oninit: oninit,
    view: function() {
        return m(
            'div',
            {
                class: 'progress',
            },
            m(
                'div', {
                    class: 'active progress-bar progress-bar-info progress-bar-striped',
                    style: 'width: ' + this.progress + '%',
                },
                '   ' + this.progress + '%   '
            )
        );
    },
};

m.mount(root, component);
