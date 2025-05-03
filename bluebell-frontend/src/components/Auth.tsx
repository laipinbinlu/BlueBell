import React, { useState } from 'react';
import { Form, Input, Button, Tabs, message } from 'antd';
import type { TabsProps } from 'antd';
import type { FormInstance } from 'antd/es/form';
import type { Rule } from 'antd/es/form';
import { userApi } from '../services/api';
import type { ApiResponse } from '../types';

interface LoginFormValues {
  username: string;
  password: string;
}

interface RegisterFormValues {
  username: string;
  password: string;
  re_password: string;
}

const Auth: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [loginForm] = Form.useForm<LoginFormValues>();
  const [registerForm] = Form.useForm<RegisterFormValues>();
  const [activeTab, setActiveTab] = useState<string>('login');

  const handleLoginSubmit = async (values: LoginFormValues) => {
    try {
      setLoading(true);
      const { username, password } = values;
      const response = await userApi.login(username, password);
      const responseData = response.data;
      
      console.log('登录响应:', responseData); // 调试日志
      
      if (responseData.code === 1000) { // 后端成功的状态码是1000
        localStorage.setItem('token', responseData.data.token);
        message.success('登录成功！');
        window.location.reload();
      } else {
        message.error(responseData.msg || '登录失败');
      }
    } catch (error: any) {
      console.error('登录失败:', error);
      message.error(error.response?.data?.msg || error.message || '登录失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleRegisterSubmit = async (values: RegisterFormValues) => {
    try {
      setLoading(true);
      const { username, password, re_password } = values;
      
      const response = await userApi.signup(username, password, re_password);
      const responseData = response.data;

      console.log('注册响应:', responseData); // 调试日志
      
      if (responseData.code === 1000) { // 后端成功的状态码是1000
        message.success('注册成功！请登录');
        setActiveTab('login');
        registerForm.resetFields();
      } else {
        message.error(responseData.msg || '注册失败');
      }
    } catch (error: any) {
      console.error('注册失败:', error);
      message.error(error.response?.data?.msg || error.message || '注册失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const validateConfirmPassword: Rule = ({ getFieldValue }) => ({
    validator(_, value) {
      if (!value || getFieldValue('password') === value) {
        return Promise.resolve();
      }
      return Promise.reject(new Error('两次输入的密码不一致'));
    },
  });

  const handleTabChange = (key: string) => {
    setActiveTab(key);
    // 切换标签页时重置表单
    if (key === 'login') {
      loginForm.resetFields();
    } else {
      registerForm.resetFields();
    }
  };

  const items: TabsProps['items'] = [
    {
      key: 'login',
      label: '登录',
      children: (
        <Form
          form={loginForm}
          layout="vertical"
          onFinish={handleLoginSubmit}
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              登录
            </Button>
          </Form.Item>
        </Form>
      ),
    },
    {
      key: 'register',
      label: '注册',
      children: (
        <Form
          form={registerForm}
          layout="vertical"
          onFinish={handleRegisterSubmit}
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>

          <Form.Item
            name="re_password"
            label="确认密码"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认密码' },
              validateConfirmPassword,
            ]}
          >
            <Input.Password placeholder="请确认密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              注册
            </Button>
          </Form.Item>
        </Form>
      ),
    },
  ];

  return (
    <div style={{ maxWidth: 400, margin: '0 auto', padding: '24px' }}>
      <Tabs activeKey={activeTab} onChange={handleTabChange} items={items} />
    </div>
  );
};

export default Auth; 