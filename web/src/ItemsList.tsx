import { useState, useEffect } from 'react';
import axios from 'axios';
import { Breadcrumb, Button, Upload, message } from 'antd';

interface ItemsListProps {
  path: string
}

function ItemsList({ path }: ItemsListProps) {
  const [items, setItems] = useState<any[]>([]);

  useEffect(() => {
    axios.get(`${process.env.REACT_APP_SERVER_URL}/api/items/${path}`)
      .then(response => {
        setItems(response.data['items'])
      })
      .catch(err => console.log(err));
  }, [])

  return (
    <div>
      <Breadcrumb items={path.split('/').map(e => { return { title: e }; })}>
      </Breadcrumb>
      <div>
        <Button
         onClick={() => {
          axios.post(`${process.env.REACT_APP_SERVER_URL}/api/items/${path}`, {
            "item_name": "asdf",
          },
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
          ).catch(err => console.log(err));
        }}>+</Button>
        <Upload
          name="items.csv"
          action={`${process.env.REACT_APP_SERVER_URL}/api/upload`}
          onChange={info => {
            if (info.file.status === "done") {
              message.success(info.file.response);
            } else if (info.file.status === "error") {
              console.log("Error uploading file")
            }
          }}>
          <Button>Upload</Button>
        </Upload>
        {items.map(e => <p key={e['item_id']}>{e['item_name']}</p>)}
      </div>
    </div>
  );
}

export { ItemsList as default };