var root = document.querySelector('#root');

var component = {
    progress: 0,
    process: function() {
        this.progress += 1;
        m.redraw()
        if (this.progress < 100) {
            var that = this;
            setTimeout(function() {
                that.process();
            }, 250);
        }
    },
    oninit: function() {
        this.process();
    },
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
