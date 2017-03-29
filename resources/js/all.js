var url = document.querySelector('body').getAttribute('data-url');

var root = document.querySelector('#root');

var oninit = function() {
    component.progress = component.progress + 1;
    m.redraw();
    if (component.progress < 100) {
        setTimeout(function() {
            oninit();
        }, 100);
    }
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
