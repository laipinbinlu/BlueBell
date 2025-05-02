import React, { useState, useEffect } from 'react';
import { Layout, Button, Modal } from 'antd';
import { PlusOutlined, LoginOutlined, LogoutOutlined } from '@ant-design/icons';
import PostList from './components/PostList';
import CreatePost from './components/CreatePost';
import Auth from './components/Auth';

const { Header, Content } = Layout;

const App: React.FC = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [showCreatePostModal, setShowCreatePostModal] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
  };

  return (
    <Layout className="layout" style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div style={{ color: 'white', fontSize: '20px' }}>BlueBell 论坛</div>
        <div>
          {isLoggedIn ? (
            <>
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => setShowCreatePostModal(true)}
                style={{ marginRight: '10px' }}
              >
                发帖
              </Button>
              <Button
                icon={<LogoutOutlined />}
                onClick={handleLogout}
              >
                退出
              </Button>
            </>
          ) : (
            <Button
              type="primary"
              icon={<LoginOutlined />}
              onClick={() => setShowAuthModal(true)}
            >
              登录/注册
            </Button>
          )}
        </div>
      </Header>

      <Content style={{ padding: '24px', backgroundColor: '#fff' }}>
        <PostList />
      </Content>

      <Modal
        title="登录/注册"
        open={showAuthModal}
        onCancel={() => setShowAuthModal(false)}
        footer={null}
        destroyOnClose
      >
        <Auth />
      </Modal>

      <Modal
        title="发布新帖子"
        open={showCreatePostModal}
        onCancel={() => setShowCreatePostModal(false)}
        footer={null}
        destroyOnClose
        width={800}
      >
        <CreatePost />
      </Modal>
    </Layout>
  );
};

export default App; 