import React from 'react';
import logo from './logo.svg'
import './App.css';
import { Layout } from 'antd';
import ItemsList from './ItemsList';

const { Header, Content, Footer } = Layout;

function App() {
  return (
    <Layout style={{ display: 'flex', height: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center' }}>
        <img src={logo} style={{ height: 50, aspectRatio: 1 }} alt='logo' />
        Header
      </Header>
      <Content>
        <ItemsList path="Home/Doujin/Hololive" />
      </Content>
      <Footer style={{ textAlign: 'right' }}>Footer</Footer>
    </Layout>
  );
}



export default App;
