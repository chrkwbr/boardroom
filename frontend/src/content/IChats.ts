export interface IChat {
  id: number;
  name: string;
  image: string;
  message: string;
  date: Date
}

export const fetchChats: () => Promise<IChat[]> = async () => {
  return initData
};

const initData: IChat[] = [
  {
    id: 0,
    name: 'Dio Lupa',
    image: 'https://img.daisyui.com/images/profile/demo/1@94.webp',
    message: '"Remaining Reason" became an instant hit, praised for its haunting sound and emotional depth. A viral\n' +
      'performance brought it widespread recognition, making it one of Dio Lupa’s most iconic tracks.',
    date: new Date()
  }, {
    id: 1,
    name: 'Ellie Beilish',
    image: 'https://img.daisyui.com/images/profile/demo/4@94.webp',
    message: '"Bears of a Fever" captivated audiences with its intense energy and mysterious lyrics. Its popularity\n' +
      'skyrocketed after fans shared it widely online, earning Ellie critical acclaim.',
    date: new Date()
  }, {
    id: 2,
    name: 'Sabrino Gardener',
    image: 'https://img.daisyui.com/images/profile/demo/3@94.webp',
    message: '"Cappuccino" quickly gained attention for its smooth melody and relatable themes. The song’s success propelled\n' +
      'Sabrino into the spotlight, solidifying their status as a rising star.',
    date: new Date()
  }
]
