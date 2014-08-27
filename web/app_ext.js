Ext.application({
    name : 'MyApp',

    launch: function(){
        Ext.widget({
            renderTo: Ext.getBody(),
            xtype: 'grid',
            title: 'Grid',
            width: '650',
            height: '300',
            plugins: 'rowediting',
            store: {
                fields: ['name', 'age', 'votes', 'credits'],
                data : [
                    [ 'Bill', 35, 10, 427],
                    [ 'Fred', 22, 4, 42]
                ]
            },
            columns: {
                defaults: {
                    editor : 'numberfield',
                    width: 120
                },
                items: [
                    {text: 'Name', dataIndex: 'name', flex:1, editor: 'textfield'},
                    {text: 'Age', dataIndex: 'age'},
                    {text: 'Votes', dataIndex: 'votes'},
                    {text: 'Credits', dataIndex: 'credits'}
                ]
            }
        })
    }
});
