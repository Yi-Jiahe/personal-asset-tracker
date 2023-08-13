import { useState, useEffect } from 'react';
import axios from 'axios';
import { Breadcrumb } from 'antd';

interface ItemsListProps {
  path: string
}

function ItemsList({ path }: ItemsListProps) {
  const [items, setItems] = useState<any[]>([]);

  useEffect(() => {
    axios.get(`https://${process.env.REACT_APP_SERVER_AUTHORITY}/api/items/${path}`)
      .then(response => {
        setItems(response.data)
      })
      .catch(err => console.log(err));
  })

  return (
    <div>
      <Breadcrumb items={path.split('/').map(e => { return { title: e }; })}>
      </Breadcrumb>
      <div>
        <button
         onClick={() => {
          axios.post(`https://${process.env.REACT_APP_SERVER_AUTHORITY}/api/items/${path}`, {
            "item_name": "asdf",
          },
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
          ).catch(err => console.log(err));
        }}>+</button>
        {items.map(e => <p key={e}>e</p>)}
      </div>
    </div>
  );
}

export { ItemsList as default };