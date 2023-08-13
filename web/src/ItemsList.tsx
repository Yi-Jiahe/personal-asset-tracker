import { useState, useEffect } from 'react';
import axios from 'axios';
import { Breadcrumb } from 'antd';

interface ItemsListProps {
  path: string
}

function ItemsList({ path }: ItemsListProps) {
  const [items, setItems] = useState<any[]>([]);

  useEffect(() => {
    axios.get('https://localhost:8080/api/items/' + path)
      .then(response => {
        setItems(response.data)
      })
      .catch(err => console.log(err));
  })

  return (
    <div>
      <Breadcrumb items={path.split('/').map(e => { return { title: e }; })}>
      </Breadcrumb>
      <div>{items.map(e => <p key={e}>e</p>)}</div>
    </div>
  );
}

export { ItemsList as default };